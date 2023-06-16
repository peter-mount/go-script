#! ./builds/linux/amd64/bin/goscript

main() {
    test()
}

// Hello world
test() {
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

    // Check we handle int/float conversions, specifically int,float => float,float
    // If that breaks then this will be an infinite loop as i starts as an int but inc is float
    for i=0;i<5;i=i+0.5 {
        print(i," ")
        if i>10 {
            print("*** FAIL ***")
            break
        }
    }
    println()

    array := newArray()
    for i:=0;i<10;i=i+1 {
        array = append( array, i )
        println(i,array)
    }
}

increment(a) {
    return a+10
}
