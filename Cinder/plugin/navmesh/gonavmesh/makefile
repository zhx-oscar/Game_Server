cc = g++
target = ../libgonavmesh.a

all:
	$(cc) -O2 -g *.cpp ./Detour/Source/*.cpp ./Recast/Source/*.cpp -I ./Detour/Include -I ./Recast/Include -c
	ar crs $(target) *.o

.PHONY : clean
clean:
	rm -f *.o
	rm -f $(target)