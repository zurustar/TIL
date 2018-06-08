#!/usr/bin/env python

'''

 以下がインストールされている必要あり

 python3 wget lv mecab

 事前に以下を実行しておく必要あり

 wget http://www2.nict.go.jp/astrec-att/member/mutiyama/manual/PostgreSQL/je.tgz
 tar zxvf ./je.tgz
 lv -Ou8 ./je/para.txt > ./para.utf8.txt

'''

import unicodedata
import subprocess
import numpy as np
import os
import keras
from keras.layers.embeddings import Embedding
from keras.layers import Dense
from keras.layers import LSTM

def run_mecab(filename):

    '''

    外部コマンド mecab を使ってfilenameで与えられたファイル内の文字列を分割（分かち書き）

    '''

    return subprocess.run(["mecab", "-Owakati", "-b65535", filename],
                            stdout=subprocess.PIPE)


def load_for_nict():

	'''

	NICTが公開している翻訳データを解析し、英語とそれに対応する日本語を抽出する
	同じフォルダに事前にUTF-8化した翻訳データがpara.utf8.txtという名前でおいておく必要あり

	'''

	ja_utf8_txt = './ja.utf8.txt'
	en_utf8_txt = './en.utf8.txt'
	ja_fp = open(ja_utf8_txt, 'w')
	en_fp = open(en_utf8_txt, 'w')
	for line in open('./para.utf8.txt').read().split('\n'):
		sentences = line.split(" ||| ")
		if len(sentences) == 3:
			ja = sentences[1]
			en = sentences[2]
			flag = True
			for c in en:
				for k in ["CJK UNIFIED", "HIRAGANA", "KATAKANA"]:
					if k in unicodedata.name(c):
						flag = False
						break
			if flag:
				en_fp.write(en + "\n")
				ja_fp.write(ja + "\n")
	ja_fp.close()
	en_fp.close()
	ja = run_mecab(ja_utf8_txt)
	en = run_mecab(en_utf8_txt)
	return en.stdout.decode(), ja.stdout.decode()


def generate_vocabulary(src):

	'''

	 文字列を受け取り、出現する文字に一意の番号を付与したうえで、
	  受け取った文字列を番号の並びに変換する

	'''
    # PADding, UNKnown, GO, End Of Sentence
	special_characters = ['<PAD>', '<UNK>', '<GO>', '<EOS>']
	vocabulary = {}
	for i in range(len(special_characters)):
		vocabulary[special_characters[i]] = i
	# 単語の最大長の取得
	max_word_length = 0
	for sentence in src.split('\n'):
		for word in sentence.split(' '):
			if max_word_length < len(word):
				max_word_length = len(word)
	array = []
	for sentence in src.split('\n'):
		sentence_array = []
		for word in sentence.split(' '):
			word_array = np.zeros(max_word_length + 1)
			for i in range(len(word)):
				character = word[i]
				if character not in vocabulary:
					vocabulary[character] = len(vocabulary)
				word_array[i] = vocabulary[character]
			sentence_array.append(word_array)
		array.append(sentence_array)
	return vocabulary, np.array(array)



def main():

	# 学習対象の翻訳データを抜き出して、
	# 学習用の形式に変換
	en, ja = load_for_nict()
	en_vocab, en_array = generate_vocabulary(en)
	ja_vocab, ja_array = generate_vocabulary(ja)

	print('en_array.shape=', en_array.shape)
	print('ja_array.shape=', ja_array.shape)


	# Embeddingの結果、何次元のテンソルにしたいか？
	output_dim = 64

	# ？？
	input_length = 10

	model = keras.models.Sequential()
	#model.add(Embedding(max_word_index, output_dim, input_length=input_length))
	model.add(Embedding(len(en_vocab), output_dim))
	model.add(LSTM(output_dim, return_sequences=True))
	model.add(LSTM(output_dim))
	#model.add(Dense(len(ja_vocab)))
	model.add(Dense(1))
	model.summary()
	model.compile(optimizer='rmsprop', loss='sparse_categorical_crossentropy')
	history = model.fit(en_array, ja_array, epochs=10, batch_size=128, validation_split=0.2)


if __name__ == '__main__':
	main()
