"""Plots benchmarking results of Arche vs. AoS and AoP"""
import numpy as np
import pandas as pd
from matplotlib import pyplot as plt

if __name__ == "__main__":
    data = pd.read_csv("results.csv", sep=";")

    model_names = {
        "Arche": "Arche",
        "ArrOfPointers": "AoP",
        "ArrOfStructs": "AoS",
        "LinkedList": "LL",
    }

    data["Model"] = ""
    data["Entities"] = 0
    data["Bytes"] = 0

    data["Benchmark"] = data["Benchmark"].str.replace("Benchmark", "")
    for index, row in data.iterrows():
        parts = row["Benchmark"].split("_")
        data.loc[index, "Model"] = model_names[parts[0]]
        data.loc[index, "Entities"] = int(parts[2]) * 1000
        data.loc[index, "Bytes"] = int(parts[1].replace("B", ""))

    data["Time"] = data["TotalTime"] / data["Entities"]

    data = data[data["Model"] != "LL"]

    models = np.unique(data["Model"])
    entities = np.unique(data["Entities"])
    bts = np.unique(data["Bytes"])

    colors = {
        "Arche": "black",
        "AoS": "red",
        "AoP": "blue",
        "LL": "grey",
    }
    linesEntities = {
        1000: ("dotted", 1.0),
        10000: ("dashed", 1.2),
        100000: ("solid", 1.5),
        # 250000: ("-.", 1.0),
    }

    plt.rcParams["svg.fonttype"] = "none"
    plt.rcParams["font.family"] = "Arial"

    fig, ax = plt.subplots(figsize=(6, 4))
    ax.set_title("Iter & get 16 byte")
    ax.set_xscale("log")
    ax.set_xticks([16, 32, 64, 128, 256])
    ax.set_xticklabels([16, 32, 64, 128, 256])
    ax.set_xlabel("Memory per Entity [byte]", fontsize=11)
    ax.set_ylabel("Time per Entity [ns]", fontsize=11)

    ax.set_ylim(0, np.max(data["Time"]) * 1.05)

    for model in reversed(models):
        for ent in entities:
            extr = data[(data["Model"] == model) & (data["Entities"] == ent)]
            extr = extr.groupby("Bytes").mean()

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
