#!/usr/bin/env python

__author__ = "@zurustar"
__status__ = "development"
__version__ = "0.0.0"
__date__ = "June 24, 2018"

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
	"""	
	各行が' ||| 'で区切られた翻訳元の英語と翻訳後の日本語の
	対になっているので、分割して英語のリストと日本語のリストにする。
	そのあとの処理の都合で、日本語には文節？ごとに半角スペースを挿入
	"""	
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
	"""
	pickleファイルがあったらそこから読み込む。
	なかったら元のテキストファイルを読み込んでpickleファイルに保存しておく。
	"""
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
	"""
	学習用データへの変換
	"""
	tokenizer = Tokenizer()
	tokenizer.fit_on_texts(texts)
	seqs = tokenizer.texts_to_sequences(texts)
	maxlen = 0
	for seq in seqs:
		if maxlen < len(seq):
			maxlen = len(seq)
	train = pad_sequences(seqs, maxlen=maxlen + 1)
	print(train.shape, train)
	return train

def main():
	# 読み込み
	en_texts, ja_texts = load()
	# 入力データの整形
	tokenizer = Tokenizer()
	tokenizer.fit_on_texts(en_texts)
	x_train = tokenizer.texts_to_sequences(en_texts)
	x_train = np.array(x_train)
	# 出力データの整形。この整形がたぶん間違ってる
	y_train = to_train(ja_texts)
	# モデル作成
	model = Sequential()
	model.add(Embedding(2, 64))
	model.add(LSTM(512, return_sequences=True))
	model.add(LSTM(512))
	model.add(Dense(2655))
	model.summary()
	model.compile(optimizer='rmsprop', loss='sparse_categorical_crossentropy')
	# 学習
	history = model.fit(x_train, y_train)

if __name__ == '__main__':
	main()

