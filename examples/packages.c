#! ./builds/linux/amd64/bin/goscript

// Hello world
main() {
    println("Packages - Math")
    println( math.Pi )
    println( math.SqrtPi )
    println( math.Min(10,5) )
    println( math.Max(10,5) )

    println( testStruct.A.C, testStruct.A.TestCall())

    println("Get c from obj.func")
    c:= testStruct.TestCall()

    println( "Get C by calling func on struct")
    println(c,c.C, c.TestCall())

    println( "Get call func then a func on it's result")
    println( testStruct.TestCall().TestCall() )

}
