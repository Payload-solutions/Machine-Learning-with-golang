<strong>KNN asumptions and pitfalls</strong>

Due to its sumplicity, KNN does not have too many assumptions,
Howeverm there some common pitfalls that you should be aware of
as you apply KNN

KNN is evaluated lazily. By this, we mean that
the distances or similarities are calclated when
need to make a prediction. There is not really anything to traing or fit prior to making a prediction. This has some advantagesm but
the calculation and search over points can be slow when you have
many data points.

The choice of k is up to you, but you should put some formalism arround chossing k and prove justification for the k that you choose. A comon technique to choose k is just to search over a eange of k values. You should, for example, start with k= 2. Then, you could start increasing k, and for each k, evaluate on a test set.


KNN doesn't take into cosinsiderations which features are more
important than other features. Moreover, if the scale of certainty
of your features is much larger than other featuresm this could unnaturally weight the importance of those larger features.
