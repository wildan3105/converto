# Load testing

## Script
```sh
#!/bin/bash

formats=(.step .iges .stl .obj)

for i in {1..100}
do
    random_format=${formats[$RANDOM % ${#formats[@]}]}

    curl --location 'http://localhost:3000/api/v1/conversions' \
    --form 'file=@"/mnt/c/Users/62823/Documents/randomfiles/100mb.shapr"' \
    --form "target_format=\"$random_format\""

    # Delay of 0.5 seconds
    sleep 0.5

done

echo "Completed 100 requests."
```