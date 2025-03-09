# Load testing

## Script
```sh
#!/bin/bash

formats=(.step .iges .stl .obj)

for i in {1..100}
do
    random_format=${formats[$RANDOM % ${#formats[@]}]}

    curl --location 'http://localhost:3000/api/v1/conversions' \
    --header 'Cookie: x-access-token=v3___eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI2Njc0ZjAzYWQyY2U5MDA4NTY0NTM0ZjgiLCJtZXJjaGFudCI6eyJfaWQiOiI2Njc0ZjAzYWQyY2U5MDA4NTY0NTM0ZjYifSwibG9naW5JZCI6IjNGUkVRRFk3MFZES0hEVjFKTDFFIiwic2lnbmVkSWQiOiI1R09JTERUUGFoS1BYNnAyZXJzQiIsInZlcnNpb24iOjMsImVtcGxveWVlIjp7Il9pZCI6IjY2NzRmMDNhZDJjZTkwMDg1NjQ1MzRmOCJ9LCJzZXNzaW9uTGVuZ3RoIjoxODAwLCJpYXQiOjE3MjI0MDc2MzEsImV4cCI6MTcyMjQ5NDAzMSwiaXNzIjoiYXV0aC52b3VjaC5zZyJ9.CKbxfw_q4UBWNjfLe6XtOf438WCuTyN4iRxEFZ_9ZE63QYbjL-t0Gl6pHlhtLQ456cYWD1OPtM9edQtpOsCLCyWYCxYqVvA8phBW_Z5OrdE1TrpiSGHC-3JlGnnu_3MW_1_sY6iMMZU75gIrzetDCER9w1aG-Im1mi4hMA3b93T2WCzQG7AAfCGEYI9hU0GasUNQwjXFltCU4Y-Uu8W0LlObABppAYawYDAn9Bkn-bTIkFNUQf5nA8jrQ_eatkRmCstl3Gi6Mt-7J2iaflZxJoc1e7SYfT_RhjOjWl00PSmN6Q-mkz-RNixUCQnzKIp5ZCHWUIE8qu1ioUuPdUd6WwVD8o1ll5IK4Q-vJoBG0Z9jWrdSMjffZqHkzyznRPbJ2jxSOUdSXH91TTIiHMrgZ10t82wC0BDeGGcBbjDdGWvZQB9ltQrYH6vVfQDN26Y8Ix0HKpcr2tC4zbjaKv40wkg1PGIQn0Yv9xpO1Zi4jWj5fojsbGZHfX_BC0A8T86-9YNtC43jDWtQhloM5FQ5D6jpFN9YWlutWInJtEKgbQkArbMXb_iXidPoulcaxlkzOxSr6im8X90saZ6pzIlRzwjOeUeXVmoETrXvY3hiui6PiiQLDd9HZRuu7QqMJOmTuG4sMaAERjrFzcx5F-T0DFlNRmSJXB1YIY8gJaYL0BQ' \
    --form 'file=@"/mnt/c/Users/62823/Documents/randomfiles/100mb.shapr"' \
    --form "target_format=\"$random_format\""

    # Delay of 0.5 seconds
    sleep 0.5

done

echo "Completed 100 requests."
```