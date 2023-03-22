import sys

if __name__ == "__main__":
    filename = sys.argv[1]
    print(filename)

    f = open(filename, "r")
    lines = f.readlines()

    for line in lines:
        print(line, end="")

    last = lines[-1].strip()
    totalStr = last.split("\t")[-1].replace("%", "")
    total = float(totalStr)

    print("\nTotal coverage: %.1f" % (total))

    if total == 100.0:
        print("Success!")
        sys.exit(0)

    print("Failed!")
    sys.exit(1)
