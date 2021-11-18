def foo(s: str) -> int:
    return len(s)


def main(args: list[str]):
    print(sum(foo(arg) for arg in args))


if __name__ == "__main__":
    import sys

    main(sys.argv[1:])
