#! ./builds/linux/amd64/bin/goscript

// Hello world
main() {
    // Reading tests
    testMap := map(
        "a": 12,
        "b": map(
            "c": 3,
            "d": map(
                "C": "wibble"
            )
        )
    )
    println( "testMap" )
    println( testMap.a )
    println( testMap.b.c )
    println( testMap.b.d.C )

    println( "testStruct")
    println( testStruct.B.C )

    println( "testSlice")
    println(testSlice)
    println(testSlice[0])
    println(testSlice[2])

    println( "testSlice2")
    println(testSlice2)
    println(testSlice2[0].C)
    println(testSlice2[1].A.C)
    println(testSlice2[1].C)
    println(testSlice2[2].C)

    for i=0;i<len(testSlice2);i=i+1 {
        println(i,testSlice2[i].C)
    }

    for i=0;i<len(testSlice);i=i+1 {
        println(i,testSlice[i])
    }

    for _,e = range testSlice {
        println(e)
    }

    println(testSlice2["2"].C)

    println(testSlice2[1+1].C)

    println(len(testSlice), len(testSlice2), len("testStruct"), len(testMap))
}
