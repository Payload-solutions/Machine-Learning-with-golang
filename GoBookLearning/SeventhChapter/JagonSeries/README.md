Understanding time series Jargon

time, chapter or timestamp: This property is the temporal element
        of each pairing in our time series. This cluld be simple a time or it could be a combination of date and time(sometimes referred to as a datetime or timestamp). It might also include time zone.

Observation, measurement, signal or random variable:This is the
        property that we are trying to forecast and/or otherwise analyze as a function of time.

Seasonality: A time series, such as the time series of airpass
        enger data, may exhibit changes that correspond to seasons(weeks, months years, and so on). time series that behave in this manner are said to exhibit some seasonality

Trends: Time series that gradually increase or decrease or de-
        crease over time (separe from seasonal effects) are said to exhibit a trend.

stationary:  A time series that exhibits the same patterns over
        time, without trends or other gradual changes (such as changes in variance or covariance), is said to be stationary

Time period: The amount of time between successive observations
        in the time series, or the difference between one timestamp and the previously ocurring timestamp in the series

Auto-regressive model: This is a model that tries to model a time
        series process by on more delayed, or lagged, versions of the same process. for example, an auto.regressive model of stokc prices would try to model stock prices by the value of the stoc price at previous time intervals

Moving average model: This is a model that tries to model a times
        series based on the current and various past values of an imperfectly predictable term, commonly referred to as error. For example, this imperfectly preditable ter may be some white noise in the time series.