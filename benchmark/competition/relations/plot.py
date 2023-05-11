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
    
    values = {
        "Entities": np.unique(data["Entities"]),
        "Parents": np.unique(data["Parents"]),
        "Children": np.unique(data["Children"]),
    }

    colors = {
        "Default": "blue",
        "Cached": "black",
        "Child": "green",
        "ParentSlice": "red",
        "ParentList": "purple",
    }
    lines_entities = {
        1000: ("dotted", 1.0),
        10000: ("dashed", 1.2),
        100000: ("solid", 1.5),
        1000000: ("solid", 2.5),
    }
    lines_parents = {
        10000: ("dotted", 1.0),
        1000: ("dashed", 1.2),
        100: ("solid", 1.5),
        10: ("solid", 2.5),
    }
    markers = {
        10: ("v", 25),
        100: ("^", 25),
        1000: ("D", 20),
        10000: ("s", 20),
        100000: ("o", 30),
        1000000: ("o", 50),
    }

    plt.rcParams["svg.fonttype"] = "none"
    plt.rcParams["font.family"] = "Arial"

    for line_column, x_column, marker_column, x_title, line_styles in [
        ("Entities", "Parents", "Children", "Parent entities", lines_entities),
        ("Entities", "Children", "Parents", "Children per parent", lines_entities),
        ("Parents", "Children", "Entities", "Children per parent", lines_parents),
    ]:
        fig, ax = plt.subplots(figsize=(7, 4))
        ax.set_title("Benchmarks of ways to represent entity relations")
        ax.set_xscale("log")
        ax.set_yscale("log")
        ax.set_xticks([10, 100, 1000, 10000])
        ax.set_xticklabels([10, 100, 1000, 10000])
        ax.set_yticks([1, 2, 5, 10, 20, 50, 100, 200])
        ax.set_yticklabels([1, 2, 5, 10, 20, 50, 100, 200])
        ax.set_xlabel(x_title, fontsize=11)
        ax.set_ylabel("Time per Entity [ns]", fontsize=11)

        ax.set_ylim(1, np.max(data["Time"]) * 1.2)

        for model in models:
            mod_data = data[(data["Model"] == model)]
            count = np.unique(mod_data[line_column])
            for ent in count:
                extr = mod_data[mod_data[line_column] == ent]
                extr = extr.groupby(x_column).mean()

                line = line_styles[ent]
                ax.plot(
                    extr.index,
                    extr["Time"],
                    linestyle=line[0],
                    linewidth=line[1],
                    color=colors[model],
                    markersize=3,
                    label=model if ent == count[-2] else None,
                    zorder=1,
                )

                marker_values = np.unique(mod_data[marker_column])
                for mk in marker_values:
                    extr2 = extr[extr[marker_column] == mk]
                    m = markers[mk]
                    ax.scatter(
                        extr2.index,
                        extr2["Time"],
                        s=m[1],
                        facecolor="white",
                        edgecolor=colors[model],
                        marker=m[0],
                        zorder=10,
                    )

        for ent in reversed(values[line_column]):
            line = line_styles[ent]
            ax.plot(
                [0],
                [0],
                linestyle=line[0],
                linewidth=line[1],
                color="black",
                label=f"{ent//1000}k {line_column}",
            )
        for ent in reversed(values[marker_column]):
            m = markers[ent]
            ax.scatter(
                [0],
                [0],
                s=m[1],
                facecolor="white",
                edgecolor=colors[model],
                marker=m[0],
                label=f"{ent} {marker_column}",
            )

        leg = ax.legend(
            bbox_to_anchor=(1.02, 1.0),
            loc="upper left",
            borderaxespad=0.0,
            fontsize=10,
        )
        leg.set(zorder=20)

        fig.tight_layout()
        fig.savefig(f"results-{line_column}{x_column}.svg")

    plt.show()
