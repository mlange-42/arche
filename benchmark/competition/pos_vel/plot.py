"""Plots benchmarking results of Arche vs. AoS and AoP"""
import numpy as np
import pandas as pd
from matplotlib import pyplot as plt

iter_prefix = "BenchmarkIter"
build_prefix = "BenchmarkBuild"

if __name__ == "__main__":
    data = pd.read_csv("results.csv", sep=";")

    models = [
        ("Arche", ("Arche\n(IDs)", "Arche\n")),
        ("ArcheGeneric", ("Arche\n(generic)", "Arche\n(batch)")),
        ("GGEcs", ("go-\ngameengine-\necs", "go-\ngameengine-\necs")),
        ("Donburi", ("Donburi", "Donburi")),
        ("Entitas", ("Entitas-Go", "Entitas-Go")),
        ("Ento", ("Ento", "Ento")),
    ]

    plt.rcParams["svg.fonttype"] = "none"
    plt.rcParams["font.family"] = "Arial"

    fig, axes = plt.subplots(ncols=2, figsize=(10, 4))

    iter = (
        axes[0],
        iter_prefix,
        [1, 2, 5, 10, 20, 50, 100, 200],
        1000,
        "log",
        "Time per iteration [μs]",
        "μs",
    )
    build = (axes[1], build_prefix, None, 1000000, "linear", "Time per run [ms]", "ms")

    for stat, (ax, prefix, yticks, factor, scale, title, units) in enumerate(
        [iter, build]
    ):
        ax.set_title("Position/Velocity -- " + prefix.replace("Benchmark", ""))
        ax.set_ylabel(title, fontsize=11)

        bench_data = data[data["Model"].str.startswith(prefix)]

        series = []
        for i, model in enumerate(models):
            extr = bench_data[
                bench_data["Model"].str.startswith(prefix + model[0] + "-")
            ]
            series.append(extr["Time"] / factor)
            median = extr["Time"].median() / factor
            ax.text(
                i + 1,
                median,
                f"{median:.1f} {units}\n",
                ha="center",
                va="baseline",
                color="#505050",
                fontsize=9,
            )

        ax.violinplot(
            series, vert=True, showmeans=True, showextrema=False, bw_method=0.5
        )

        ax.set_xticks(range(1, len(models) + 1))
        ax.set_xticklabels([n[stat] for _, n in models])

        ax.set_ylim((yticks or [0])[0], np.max(bench_data["Time"]) * 1.1 / factor)
        ax.set_yscale(scale)

        if yticks is not None:
            ax.set_yticks(yticks)
            ax.set_yticklabels(yticks)

    fig.tight_layout()
    fig.savefig("results.svg")
    plt.show()
