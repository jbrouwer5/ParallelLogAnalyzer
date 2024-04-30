#!/bin/bash
#
#SBATCH --mail-user=jbrouwer@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj3_benchmark
#SBATCH --chdir=/home/jbrouwer/parallel/project-3-jbrouwer5/proj3/
#SBATCH --output=./slurm/out/performance.stdout
#SBATCH --error=./slurm/out/error.stderr
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --time=5:00


module load golang/1.19
go run logProcessor/logProcessor.go < input.txt 
go run logProcessor/logProcessor.go 2 < input.txt 
go run logProcessor/logProcessor.go 4 < input.txt 
go run logProcessor/logProcessor.go 6 < input.txt 
go run logProcessor/logProcessor.go 8 < input.txt 
go run logProcessor/logProcessor.go 12 < input.txt 

go run logProcessor/logProcessor.go < largeInput.txt 
go run logProcessor/logProcessor.go 2 < largeInput.txt 
go run logProcessor/logProcessor.go 4 < largeInput.txt 
go run logProcessor/logProcessor.go 6 < largeInput.txt 
go run logProcessor/logProcessor.go 8 < largeInput.txt 
go run logProcessor/logProcessor.go 12 < largeInput.txt 

python3 graph.py 