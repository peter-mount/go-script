#! ./builds/linux/amd64/bin/goscript

// Hello world
main() {
    // Reading tests
    println( "testMap" )
    println( testMap.a )
    println( testMap.b.c )
    println( testMap.b.d.C )

    println( "testStruct")
    println( testStruct.B.C )
}
