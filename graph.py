import re
import matplotlib.pyplot as plt

# Initialize lists to hold the data
small_input_times = []
large_input_times = []
threads = [1, 2, 4, 6, 8, 12]  
time1 = 0 
time2 = 0 

# Read the file and extract times
with open('slurm/out/performance.stdout', 'r') as file:
    content = file.read()
    
    # Split the content by the pattern that represents each call section
    sections = content.split("\n")
    for i in range(12):
        time_found = sections[(i+1)*4-1]
        if time_found:
            time = float(time_found[:-2])
            if (i == 0):
                time1 = time
            
            if (i == 6):
                time2 = time
            
            if i < 6:  # First 6 for small inputs
                small_input_times.append(time1/time)
            else:  # Last 6 for large inputs
                large_input_times.append(time2/time)

# Plotting
plt.figure(figsize=(10, 6))
plt.plot(threads, small_input_times, marker='o', linestyle='-', color='b', label='Small Inputs')
plt.plot(threads, large_input_times, marker='o', linestyle='-', color='r', label='Large Inputs')
plt.title('Speedup vs Number of Threads')
plt.xlabel('Number of Threads')
plt.ylabel('Speedup')
plt.legend()
plt.grid(True)
plt.savefig('speedup_graph.png')
