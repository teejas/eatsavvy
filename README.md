# eatsavvy

uses Google Places API and Vapi to enrich restaurant details with nutritional and dietary information.

## api

accepts search query to find restaurants to enrich.

## worker

uses Vapi to call restaurants and collect information from restaurants.

# Notes

the architecture is intended to decouple initiation of the phone call from the API

this is because restaurants may not be open when the request to enrich data is submitted, in which case we
would have to poll periodically or have a long timeout or do something else heinous

we really only need one publisher (the api) and one consumer (the worker) though since all the work is
relatively short/easy

1. user submits a job (enrich nutrition info about this restaurant)
2. job is added to queue
3. worker picks up a job
3a. if restaurant is closed, it re-queues with delay
3b. if restaurant is open, it places outbound call via Vapi API
4. vapi makes phone call and sends eocr to /process-eocr with both transcript and structured outputs
5. we either process the transcript in API or also queue it for an offline worker to handle

## To-do
[ ] refactor internal/worker/*
[ ] dynamically generate structured outputs for assistant (using structuredMultiData)
[ ] add Yelp support (for reviews and supplementing missing phone numbers)