CPP = g++

CPPFLAGS  = -std=c++17
CPPFLAGS += -Wall
CPPFLAGS += -Wextra
CPPFLAGS += -pedantic
CPPFLAGS += -Werror

VFLAGS	= --quiet
VFLAGS += -v
VFLAGS += --tool=memcheck
VFLAGS += --leak-check=full
VFLAGS += --error-exitcode=1
VFLAGS += --show-reachable=yes
VFLAGS += --show-possibly-lost=yes
VFLAGS += --undef-value-errors=yes
VFLAGS += --track-origins=yes

CLI  = arduino-cli
FQBN = arduino:samd:mkr1000
PORT = /dev/ttyACM99
INO  = raingauge.ino
WARN = --warnings all
BINDIR = bin

build: 
	$(CLI) compile --fqbn $(FQBN) $(INO) $(WARN)

upload:
	$(CLI) upload -p $(PORT) --fqbn $(FQBN) $(INO) 

test_tlv.out: test/test_tlv.cpp src/tlv.cpp
	@$(CPP) $(CPPFLAGS) -o $(BINDIR)/$@ $^
	@./$(BINDIR)/$@

memcheck: test_tlv.out
	@valgrind $(VFLAGS) $(BINDIR)/$^

test: test_tlv.out

clean: 
	$(CLI) cache clean
	rm $(BINDIR)/*.out

all: build upload

.PHONY: build test
