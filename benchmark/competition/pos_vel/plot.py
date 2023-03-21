"""Plots benchmarking results of Arche vs. AoS and AoP"""
import numpy as np
import pandas as pd
from matplotlib import pyplot as plt

iter_prefix = "BenchmarkIter"

if __name__ == "__main__":
    data = pd.read_csv("results.csv", sep=";")

    models = [
        ("Arche", "Arche"),
        ("GGEcs", "go-gameengine-ecs"),
        ("Donburi", "Donburi"),
        ("Ento", "Ento"),
    ]

    plt.rcParams["svg.fonttype"] = "none"
    plt.rcParams["font.family"] = "Arial"

    fig, ax = plt.subplots(figsize=(6, 4))

    ax.set_title("Position/Velocity benchmark")
    ax.set_xlabel("Engine", fontsize=11)
    ax.set_ylabel("Time per iteration [us]", fontsize=11)

    bench_data = data[data["Model"].str.startswith(iter_prefix)]

    series = []
    for model in models:
        extr = bench_data[bench_data["Model"].str.startswith(iter_prefix + model[0])]
        series.append(extr["Time"] / 1000)
        """ax.plot(
            extr["Bytes"],
            extr["Time"],
            linestyle=line[0],
            linewidth=line[1],
            color=colors[model],
            marker="o",
            markersize=3,
            label=model if ent == 100000 else None,
        )"""
    ax.violinplot(series, showmeans=True, showextrema=False, bw_method=0.5)

    ax.set_xticks(range(1, len(models) + 1))
    ax.set_xticklabels([n for _, n in models])

    ax.set_ylim(2, np.max(bench_data["Time"]) * 1.1 / 1000)
    ax.set_yscale("log")
    ax.set_yticks([2, 5, 10, 20, 50, 100])
    ax.set_yticklabels([2, 5, 10, 20, 50, 100])

    fig.tight_layout()
    fig.savefig("results.svg")
    plt.show()
