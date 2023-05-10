"""Plots benchmarking results of Arche vs. AoS and AoP"""
import numpy as np
import pandas as pd
from matplotlib import pyplot as plt

if __name__ == "__main__":
    data = pd.read_csv("results.csv", sep=";")

    data["Model"] = ""
    data["Entities"] = 0
    data["Parents"] = 0
    data["Children"] = 0

    data["Benchmark"] = data["Benchmark"].str.replace("BenchmarkRelation", "")
    for index, row in data.iterrows():
        parts = row["Benchmark"].split("-")[0].split("_")
        data.loc[index, "Model"] = parts[0]
        data.loc[index, "Parents"] = int(parts[2])
        data.loc[index, "Children"] = int(parts[4])

    data["Entities"] = data["Parents"] * data["Children"]
    data["Time"] = data["TotalTime"] / data["Entities"]

    data = data[data["Entities"] > 1000]

    models = ["ParentList", "ParentSlice", "Child", "Default", "Cached"]
    parents = np.unique(data["Parents"])
    children = np.unique(data["Children"])

    colors = {
        "Default": "blue",
        "Cached": "black",
        "Child": "green",
        "ParentSlice": "red",
        "ParentList": "purple",
    }
    linesEntities = {
        1000: ("dotted", 1.0),
        10000: ("dashed", 1.2),
        100000: ("solid", 1.5),
        1000000: ("solid", 2.5),
    }

    plt.rcParams["svg.fonttype"] = "none"
    plt.rcParams["font.family"] = "Arial"

    fig, ax = plt.subplots(figsize=(6, 4))
    ax.set_title("Iter & get 16 byte")
    ax.set_xscale("log")
    ax.set_yscale("log")
    ax.set_xticks([10, 100, 1000, 10000])
    ax.set_xticklabels([10, 100, 1000, 10000])
    ax.set_yticks([1, 2, 5, 10, 20, 50, 100, 200])
    ax.set_yticklabels([1, 2, 5, 10, 20, 50, 100, 200])
    ax.set_xlabel("Number of parents", fontsize=11)
    ax.set_ylabel("Time per Entity [ns]", fontsize=11)

    ax.set_ylim(1, np.max(data["Time"]) * 1.2)

    for model in models:
        mod_data = data[(data["Model"] == model)]
        entities = np.unique(mod_data["Entities"])
        for ent in entities:
            extr = mod_data[mod_data["Entities"] == ent]
            extr = extr.groupby("Parents").mean()

            line = linesEntities[ent]
            ax.plot(
                extr.index,
                extr["Time"],
                linestyle=line[0],
                linewidth=line[1],
                color=colors[model],
                marker="o",
                markersize=3,
                label=model if ent == 100000 else None,
            )
    
    for ent in reversed(entities):
        line = linesEntities[ent]
        ax.plot(
            [0],
            [0],
            linestyle=line[0],
            linewidth=line[1],
            color="black",
            label=f"{ent//1000}k entities",
        )

    ax.legend()

    fig.tight_layout()
    fig.savefig("results.svg")
    plt.show()
