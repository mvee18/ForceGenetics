import os
import sys
import numpy as np

if __name__ == "__main__":
	difference = []

	f1 = os.path.abspath(sys.argv[1])
	f2 = os.path.abspath(sys.argv[2])

	a1 = np.genfromtxt(f1, skip_header=1)
	a2 = np.genfromtxt(f2, skip_header=1)

	af1 = a1.flatten()
	af2 = a2.flatten()

	af1 = np.sort(af1)
	af2 = np.sort(af2)

	difference = np.subtract(af1, af2)

	print(difference)

# Conclusion: The fort.15 files generated by the algorithm are not the same.