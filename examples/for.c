#! ./builds/linux/amd64/bin/goscript

// Hello world
main() {

    i=0
    println("i before loop",i)
    for i=0; i<len(testSlice2); i=i+1 {
        println(i,testSlice2[i].C)
    }
    println("i after loop",i)

    i=42
    println("i before loop",i)
    for i:=0; i<len(testSlice); i=i+1 {
        println(i,testSlice[i])
    }
    println("i after loop",i)

    println("before range",i)
    for i,e := range testSlice {
        println(i,e)
    }
    println("after range",i)

    println("before try")
    try (a=testcl0 ; b=testcl1 ) {
        println(a,b)
        for i:=0;i<10;i=i+1 {
            println("try",i)
            if i>5 break
        }
        println(testSlice2[7677])
    }
    catch(ex) {
        println("Exception", ex)
    }
    finally {
        println("Finally block start")
    }
    println("after try")

    m := map( "a":1, "b":2, "c": 3.1415926 )
    println(m)
    for k,v := range m { println(k,v)}
    println("c=",m.c)
    println("c=",m["c"])

}
