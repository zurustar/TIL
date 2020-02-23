
PostgreSQLを前提としたRDBMSの勉強

# RDBMS



# PostgreSQL



# SQL

RDBMSを操作するために標準化された言語。

## 結合
ふたつのテーブルをくっつける方法は２種類。
縦方向にくっつける」「UNION」と横方向にくっつける「JOIN」がある。
JOINにはさらに
「(INNER) JOIN」「LEFT (OUTER) JOIN」「RIGHT (OUTER) JOIN」「CROSS JOIN」
という結合方法が異なる４種類がある。

### UNION

|  X  |  Y  |  Z  |
| --- | --- | --- |
|  1  |  3  |  5  |
|  2  |  4  |  6  |

と

|  X  |  Y  |  Z  |
| --- | --- | --- |
|  7  |  9  |  b  |
|  8  |  a  |  c  |

を UNIONすると

|  X  |  Y  |  Z  |
| --- | --- | --- |
|  1  |  3  |  5  |
|  2  |  4  |  6  |
|  7  |  9  |  b  |
|  8  |  a  |  c  |

になる。SQLだと、

SELECT table1 UNION table2;

みたいに書く。

○ JOIN


