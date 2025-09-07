package types

/**
    Centralized Container stats struct
**/
type ContainerStats struct {
    ID string
    NAME string
    CPUUSAGE uint64
    MEMORYUSAGE uint64
    MEMORYLIMIT uint64
}
