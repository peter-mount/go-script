#! ./builds/linux/amd64/bin/goscript

// Hello world
main() {
    b = 1
    println(b)
    c = increment(b)
    println("return from func", c,b)
    c = increment(b=5)
    println("return from func", c,b)

    for i=0; i<10; i=i+1 {
        print(i)
    }
    println()

    for i=0; i<10; i=i+1 print(i)
    println()

    j=0
    for ;j<10;j=j+1 print(j)
    println()

    j=0
    for ;j<10; {
        print(j)
        j=j+1
    }
    println()

    j=0
    for ;; {
        print(j," ")
        j=j+1
        if j>9 {
            print("break")
            break
        }
    }
    println()

    // Between test
    min = 5
    max = 10
    for i=0;i<12;i=i+1 {
        println( min, "<=", i, "<=", max, between(i,min,max) )
    }
}

increment(a) {
    return a+10
}
