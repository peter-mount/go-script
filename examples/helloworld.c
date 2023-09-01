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

        ary := map( "a":1, "b":2, "c":3, "d":4 )
        for i,v:=range ary fmt.Printf("i=%q v=%d\n",i,v)
        for _,v:=range ary fmt.Printf("key _ v=%d\n",v)

        for i:=0;i<10;i=i+1 fmt.Printf("for i=%d\n")

        a:=1
        b:=10
        repeat {
          fmt.Printf("repeat %d until %d",a,b)
          a=a+1
        } until a>b

    } catch( err ) {
        println(err)
    }
}

test() {
    a = 2
    b = 84
    fmt.Printf(" %d / %d = %d\n", b, a, b/a )
}

test2() {
    return map( "a": 1, "b": 2)
}