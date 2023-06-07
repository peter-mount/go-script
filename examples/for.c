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
}
