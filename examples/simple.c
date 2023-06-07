#! ./builds/linux/amd64/bin/goscript

/* This is an example uC program. */
fac(n)
{
    if (n < 2)
        return n;
    return n * fac(n - 1);
}

sum(int n, int a[]) {
    int i;
    int s;

    i = 0;
    s = 0;
    while (i <= n) {
        s = s + a[i];
        i = i + 1;
    }
    return s;
}

main()
{
    int a[2];

    a[0] = fac(5);
    a[1] = 27;
    putint(sum(2, a)); // prints 147
    return 0;
}