NAME:=aoc2023
BIN_DIR:=./bin

PROG:=$(BIN_DIR)/$(NAME)

GOFILES:=$(shell find src/ -type f -name "*.go")

NOWDATE:=$(shell TZ="EST" date +%Y%m%d)
NOWDAY:=$(shell TZ="EST" date '+%e')
ENDDATE:=20231225
DOCKERRUN=docker run --rm -i --env AOC_SESSION ${AOC_RUNOPTS} aoc2023:latest $(ELAPSEDOPTS)
ifdef ELAPSED
ELAPSEDOPTS:=-e
endif

.PHONY: build run run-all clean example build-run run-bare example-bare all today diy-run

all: build

build: $(PROG)


$(PROG): $(GOFILES)
	DOCKER_BUILDKIT=1 docker build --target bin --output $(BIN_DIR)/ . 
	touch $(PROG)

build-run: $(PROG)
	docker build -f Dockerfile.run -t aoc2023:latest .

run: build-run $(PROG)
	@$(PROG) input $(DAY) | $(DOCKERRUN) $(DAY)

run-bare: $(PROG)
	@$(PROG) input $(DAY) | $(PROG) run $(DAY)

run-all: $(PROG)
	@if test "$(NOWDATE)" -lt "$(ENDDATE)"; then for day in `seq $(NOWDAY)`; do $(PROG) input $$day | $(DOCKERRUN) $$day; done; else for day in `seq 25`; do $(PROG) input $$day | $(DOCKERRUN) $$day;done;fi

today: build-run $(PROG)
	@$(PROG) input $(NOWDAY) | $(DOCKERRUN) $(NOWDAY)

clean:
	rm -f $(PROG)

example: $(PROG) build-run
	@cat examples/day$(DAY).txt | $(DOCKERRUN) $(DAY)

example-bare: $(PROG)
	@cat examples/day$(DAY).txt | $(PROG) run $(DAY)

diy-run: build-run $(PROG)
	$(DOCKERRUN) $(DAY)
