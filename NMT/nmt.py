#!/usr/bin/env python

#
# 以下がインストールされている必要あり
#
# python3 wget lv mecab
#
# 事前に以下を実行しておく必要あり
#
# wget http://www2.nict.go.jp/astrec-att/member/mutiyama/manual/PostgreSQL/je.tgz
# tar zxvf ./je.tgz
# lv -Ou8 ./je/para.txt > ./para.utf8.txt
#

import unicodedata
import subprocess
import os


def load():
	mecab = 'mecab'
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
	ja = subprocess.run([mecab, "-Owakati", "-b65535", ja_utf8_txt], stdout=subprocess.PIPE)
	en = subprocess.run([mecab, "-Owakati", "-b65535", en_utf8_txt], stdout=subprocess.PIPE)
	return en.stdout.decode(), ja.stdout.decode()


def generate_vocabulary(src):
	vocabulary = {}
	for sentence in src.split('\n'):
		for word in sentence.split(' '):
			for character in word:
				if character not in vocabulary:
					vocabulary[character] = len(vocabulary)
	return vocabulary

def conv2ary(en, ja):
	en_vocabulary = generate_vocabulary(en)
	ja_vocabulary = generate_vocabulary(ja)
	print(en_vocabulary)
	print(ja_vocabulary)

def main():
	en, ja = load()
	conv2ary(en, ja)

if __name__ == '__main__':
	main()

