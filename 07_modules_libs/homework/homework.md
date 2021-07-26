# Homework 04

## Stock price service variant 1

Implement http server with one endpoint `GET /price/{ticker}/date/{date}` which returns the stock price for exact ticker 
on exact date (e.g. `GET /price/AAPL/date/2021-07-22` -> `{"ticker":"AAPL","close_price":"146.80","date":"2021-07-22"}`)

To obtain stock price info you can use external API (e.g. https://www.alphavantage.co/)

## Stock price service variant 2

Implement http server with one endpoint `GET /price/{ticker}/stat` which returns the highest, lowest 
and average prices (e.g. `GET /price/AAPL/stat` -> `{"ticker":"AAPL","highest_price":"149.80","lowest_price":"34.80","avg_price":"97.56"}`)

To obtain stock price info you can use external API (e.g. https://www.alphavantage.co/)

## Stock price service variant 3

Implement http server with one endpoint `GET /price/{ticker}/diff/{first_date}/{second_date}` which returns close prices 
percentage difference (e.g. `GET /price/AAPL/diff/2021-07-15/2021-07-22` -> `{"ticker":"AAPL","percentage_diff":"-1.13","first_date":"2021-07-15","second_date":"2021-07-22"`)

To obtain stock price info you can use external API (e.g. https://www.alphavantage.co/)
