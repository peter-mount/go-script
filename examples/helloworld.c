#! ./builds/linux/amd64/bin/goscript
#include "test.c"
//include "./examples/test.c"

// Hello world
main() {
    try {
        // Integers
        println(3)
        // Floating point
        println(3.14159263)
        println(0.5)
        println(5.0)
        // Strings
        {
        println("Hello World!")
        println("Basic!")
        }

        {
            // Set a and print it
            a=3.14159263
            println("a=3.14159263", a)

            // Set a and b should be the same
            a = (b = 42)
            println("a = (b = 42)", a, b)

            a = 10
            b = 42
            println("a+b=", a+b)
            println("a-b=", a-b)
            println("a*b=", a*b)
            println("a/b=", a/b)
            a = 10.0
            println("a/b=", a/b)
            c=a+b
            println("c=", c)

            a = 10 + (b = 42)
            println(a, b)
        }

        a = 10
        b = 20
        println( "a==b", a, b, a!=2)

        a = 1
        while a < 10 {
            print(a," ")
            a=a+1
        }

        print("=",a)
        println("Success")

        try( file := os.Create("/tmp/test.txt")) {
            println( "file", file)
            fmt.Fprintln( file, "Hello\n")
        }

        file2 := os.Create("/tmp/test.txt")
        try( file2 ) {
        }

        println("**********************************")

        test()

        v:=test2()
        println( v.a )
        println( v.b )

        println("**********************************")
        println("Loop tests\n")

        ary := map( "a":1, "b":2, "c":3, "d":4 )

        fmt.Println("for k,v = range")
        c:=0
        for k,v:=range ary {
            fmt.Printf(" %q=%d",k,v)
            c=c+1
        }
        result( c==4 )

        fmt.Println("for k,v = range")
        c:=0
        for _,v:=range ary {
            fmt.Printf(" _=%d",v)
            c=c+1
        }
        result( c==4 )

        // These should just work
        testFor(0,10,1,10)
        testRepeat(1,10,10)
        testWhile(1,10,9)

        // These should not go into an infinite loop
        // as start > end. Repeat should run once but
        // while should not run

        // a>b on start but should still run once
        testRepeat(11,10,1)
        // a>b on start so should never run
        testFor(11,10,1,0)
        testWhile(11,10, 0)

    } catch( err ) {
        println(" FAIL")
        println(err)
    }
}

testFor(a, b, s, count) {
    c:=0
    fmt.Printf("for i=%d; i<%d; i=i+%d\n",a,b,s)
    for i:=a; i<b; i=i+s {
        fmt.Printf(" %d",a)
        a=a+1
        c=c+1

        // Catch if overrunning
        if (a-b)>10 throw( "for broken")
    }
    result( c == count)
}

testRepeat(a, b, count) {
    c:=0
    fmt.Printf("a=%d; repeat a=a+1 until a>%d\n",a,b)
    repeat {
        fmt.Printf(" %d",a)
        a=a+1
        c=c+1

        // Catch if overrunning
        if (a-b)>10 throw( "repeat broken")
    } until a>b
    result( c == count)
}

testWhile(a, b, count) {
    c:=0
    fmt.Printf("a=%d; while a<%d a=a+1\n",a,b)
    while a<b {
        fmt.Printf(" %d",a)
        a=a+1
        c=c+1

        // Catch if overrunning
        if (a-b)>10 throw( "while broken")
    }
    result( c == count)
}

result(c) {
    if c fmt.Println(" PASS\n") else fmt.Println(" FAIL\n")
}
test() {
    a = 2
    b = 84
    fmt.Printf(" %d / %d = %d\n", b, a, b/a )
}

test2() {
    return map( "a": 1, "b": 2)
}