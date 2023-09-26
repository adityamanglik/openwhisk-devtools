# Don't stop on errors for these commands
rm *.txt || true

source Experiment.sh http://128.110.96.69:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/hello/world JS
source Experiment.sh http://128.110.96.69:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloJava/world Java

python response_time_plotter.py JavaOutputTime.txt
python response_time_plotter.py JSOutputTime.txt