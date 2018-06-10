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
import random
import keras
from keras.models import Sequential
from keras.layers.embeddings import Embedding
from keras.layers import Input
from keras.layers import Dense
from keras.layers import LSTM

# ----------------------------------------------------------------------------
#
#
#
#
def run_mecab(filename):
	return subprocess.run(["mecab", "-Owakati", "-b65535", filename],
							stdout=subprocess.PIPE)


# ----------------------------------------------------------------------------
#
# メモリが少ない環境で遊ぶ場合はmax_linesで値を指定して元ネタを絞る
#
def load_for_nict(max_lines=None):
	lines = open('./para.utf8.txt').read().split('\n')
	# メモリが少ないとハングするのでmax_linesが指定されていたら絞る
	if max_lines != None:
		if len(lines) > max_lines:
			random.shuffle(lines)
			lines = lines[:max_lines]
	ja_utf8_txt = './ja.utf8.txt'
	en_utf8_txt = './en.utf8.txt'
	ja_fp = open(ja_utf8_txt, 'w')
	en_fp = open(en_utf8_txt, 'w')
	for line in lines:
		sentences = line.split(" ||| ")
		if len(sentences) == 3:
			ja = sentences[1]
			en = sentences[2]
			# たまに英語側に日本語が混ざっている異常データがあるので抹消
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


#
#
#
#
#
def preprocess(src):
	#
	# (1) 箱の準備
	#	 PADding, UNKnown, GO, End Of Sentence
	#
	special_characters = ['<PAD>', '<UNK>', '<GO>', '<EOS>']
	vocabulary = {}
	for i in range(len(special_characters)):
		vocabulary[special_characters[i]] = i
	#
	# (2) 文書の最大長の取得
	#
	max_sentence_length = 0
	max_word_length = 0
	for sentence in src.split('\n'):
		words = sentence.split(' ')
		if max_sentence_length < len(words):
			max_sentence_length = len(words)
			print(" 文章長を更新しました(", max_sentence_length, ")：",words)
		for word in words:
			if max_word_length < len(word):
				max_word_length = len(word)
				print(" 単語長を更新しました(", max_word_length, ")：", word)
	print("最長文書は", max_sentence_length, "単語あった")
	print("最長単語は", max_word_length, "文字あった")
	#
	# (3) 数値配列化
	#
	ary = []
	for sentence in src.split('\n'):
		sentences_ary = []
		words = sentence.split(' ')
		for word in words:
			words_ary = np.zeros(max_word_length + 1)
			for i in range(len(word)):
				character = word[i]
				# いま調べている文字が辞書になかったら追加
				if character not in vocabulary:
					vocabulary[character] = len(vocabulary)
				words_ary[i] = vocabulary[character]
			sentences_ary.append(words_ary)
		ary.append(sentences_ary)
	return vocabulary, np.array(ary)

def embedding_sample():
	vocab_size = 5
	input_array = np.random.randint(vocab_size, size=(4, 3))
	print(input_array)
	model = keras.models.Sequential()
	model.add(keras.layers.embeddings.Embedding(vocab_size, 4, input_length=3))
	model.compile('rmsprop', 'mse')
	output_array = model.predict(input_array)
	print(output_array)


# ----------------------------------------------------------------------------
#
#
#
#
def main():

	print("--- load ---")
	en, ja = load_for_nict(100)
	print('len(en) =', len(en), ' len(ja) =', len(ja))
	print()

	print("--- preprocess EN ---")
	en_vocab, en_ary = preprocess(en)
	print('len(en_ary) =', len(en_ary))
	print('len(en_ary[0]) =', len(en_ary[0]))
	print('len(en_ary[0][0]) =', len(en_ary[0][0]))
	print('en_ary[0][0] =[', en_ary[0][0], ']')
	print()

	print("--- preprocess JA ---")
	ja_vocab, ja_ary = preprocess(ja)
	print('len(ja_ary) =', len(ja_ary))
	print('len(ja_ary[0]) =', len(ja_ary[0]))
	print('len(ja_ary[0][0]) =', len(ja_ary[0][0]))
	print('ja_ary[0][0] =', ja_ary[0][0])
	print()

	print("--- create model ---")
	model = keras.models.Sequential()
	model.add(Input())
	model.add(Embedding(len(en_vocab), 64))
	model.add(LSTM(512, return_sequences=True))
	model.add(LSTM(512))
	model.add(Dense(len(ja_vocab)))
	model.summary()
	model.compile(optimizer='rmsprop', loss='sparse_categorical_crossentropy')
	print()

	print("--- train ---")
	print("en_ary.shape =", en_ary.shape)
	history = model.fit(en_ary, ja_ary, epochs=10, batch_size=32, validation_split=0.2, verbose=2)


if __name__ == '__main__':
	main()
