import matplotlib.pyplot as plt
import os 
import sys

x = []
y = []

def convert_time(t):
	minutes = t.split('m')
	seconds = minutes[1].split('s')[0]
	
	total = float(minutes[0])*60+float(seconds)

	return total

if __name__ == "__main__":
	input_file = os.path.abspath(sys.argv[1])
	file1 = open(input_file, 'r')
	lines = file1.readlines()

	for line in lines:
		if "Time taken so far: " in line:
			splitted = line.split()
			time = splitted[4]
			fitness = splitted[10]

			time_seconds = convert_time(time)

			x.append(time_seconds)
			y.append(float(fitness))

	print(x, y)

	fig, ax = plt.subplots()

	ax.plot(x, y)

	ax.set(xlabel='time (min)', ylabel='fitness',
	title='Fitness versus Time')
	ax.grid()

	fig.savefig("test.png")
	# plt.show()

