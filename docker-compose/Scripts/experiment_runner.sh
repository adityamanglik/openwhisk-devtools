# Don't stop on errors for these commands
rm JavaOutput.txt || true
rm JSOutput.txt || true

source JavaExperiment.sh
python response_time_plotter.py JavaOutput.txt

source JSExperiment.sh
python response_time_plotter.py JSOutput.txt