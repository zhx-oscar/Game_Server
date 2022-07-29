cc = g++
target_release = ./install/lib/release/libPhysXCWrap.a
target_debug = ./install/lib/debug/libPhysXCWrap.a
target_h = ./install/incl/PhysXCWrap.h
physx_home = ../PhysX-4.1/

release:
	$(cc) -std=c++11 -O2 -g -DNDEBUG *.cpp -I $(physx_home)/install/linux/PhysX/include -I $(physx_home)/install/linux/PxShared/include -c
	ar crs $(target_release) *.o
	cp -f ./PhysXCWrap.h $(target_h)

debug:
	$(cc) -std=c++11 -O0 -g -D_DEBUG *.cpp -I $(physx_home)/install/linux/PhysX/include -I $(physx_home)/install/linux/PxShared/include -c
	ar crs $(target_debug) *.o
	cp -f ./PhysXCWrap.h $(target_h)

.PHONY : clean
clean:
	rm -f *.o
	rm -f $(target_release)
	rm -f $(target_debug)
	rm -f $(target_h)