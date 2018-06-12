#!/usr/bin/env python

import os
import numpy as np
import pickle
from janome.tokenizer import Tokenizer as janome_tokenizer
from keras.preprocessing.text import Tokenizer
from keras.preprocessing.sequence import pad_sequences

from keras.models import Sequential
from keras.layers.embeddings import Embedding
from keras.layers import LSTM
from keras.layers import Dense

MAX_SEQUENCE_LENGTH = 1000

def _load(filename):
	t = janome_tokenizer()
	ja_texts, en_texts = [], []
	lines = open(filename).read().split('\n')
	for line in lines:
		elms = line.split(' ||| ')
		if len(elms) != 3:
			continue
		en_texts.append(elms[2])
		ja_list = []
		for token in t.tokenize(elms[1]):
			ja_list.append(token.surface)
		ja_texts.append(" ".join(ja_list))
		print(len(en_texts), "/", len(lines))
	return [en_texts, ja_texts]

def load():
	filename = './para.utf8.pickle'
	data = []
	if os.path.exists(filename):
		fp = open(filename, mode='rb')
		data = pickle.load(fp)
		fp.close()
	else:
		data = _load('./para.utf8.txt')
		fp = open(filename, mode='wb')
		pickle.dump(data, fp)
		fp.close()
	return data[0], data[1] # en, ja

def to_train(texts):
	tokenizer = Tokenizer()
	tokenizer.fit_on_texts(texts)
	seqs = tokenizer.texts_to_sequences(texts)
	return np.array(seqs)
	maxlen = 0
	for seq in seqs:
		if maxlen < len(seq):
			maxlen = len(seq)
	train = pad_sequences(seqs, maxlen=maxlen + 1)
	print(train.shape, train)
	return train

def main():
	en_texts, ja_texts = load()
	tokenizer = Tokenizer()
	tokenizer.fit_on_texts(en_texts)
	x_train = tokenizer.texts_to_sequences(en_texts)
	x_train = np.array(x_train)
	y_train = to_train(ja_texts)
	model = Sequential()
	model.add(Embedding(2, 64))
	model.add(LSTM(512, return_sequences=True))
	model.add(LSTM(512))
	model.add(Dense(2655))
	model.summary()
	model.compile(optimizer='rmsprop', loss='sparse_categorical_crossentropy')
	history = model.fit(x_train, y_train)

if __name__ == '__main__':
	main()

