def fuel(mass):
    return (mass // 3) - 2


def rec_fuel(mass):
    mass = fuel(mass)
    total = 0
    while mass > 0:
        total += mass
        mass = fuel(mass)
    return total


def main(args):
    total_fuel = 0
    total_rec_fuel = 0
    with open(args[1]) as f:
        for line in f:
            total_fuel += fuel(int(line))
            total_rec_fuel += rec_fuel(int(line))
    print(total_fuel)
    print(total_rec_fuel)


if __name__ == '__main__':
    import sys
    main(sys.argv)
