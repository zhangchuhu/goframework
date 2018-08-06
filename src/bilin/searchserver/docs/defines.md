# 搜索类型

## user

使用 `typ=1` 搜索，使用 topic `bilin_user_update` 或 `bilin_user_update_test` 更新

字段名（区分大小写）, 字段的取值，类型统一，都是 string

| 字段名        | 字段的意义                | 是否加入搜索 | 是否外部存储 |
| ------------- | ------------------------- | ------------ | ------------ |
| id            | user id                   |              |              |
| bilin_id      | 比邻号                    | Y            |              |
| name          | 昵称                      | Y            |              |
| avatar        | 头像                      |              | Y            |
| sex           | 性别                      |              | Y            |
| age           | 年龄                      |              | Y            |
| location      | 位置                      |              | Y            |
| live          | `0` 没开播， `1` 正在直播 |              | Y            |
| room_user_num | 如果正在直播，为房间人数  |              | Y            |

## room

使用 `typ=2` 搜索，使用 topic `bilin_room_update` 或 `bilin_room_update_test` 更新

字段名（区分大小写）, 字段的取值，类型统一，都是 string

| 字段名     | 字段的意义                                        | 是否加入搜索 | 是否外部存储 |
| ---------- | ------------------------------------------------- | ------------ | ------------ |
| id         | room id                                           |              |              |
| name       | 房间标题                                          | Y            |              |
| live       | `0` 没开播， `1` 正在直播                         | Y            |              |
| display_id | 官频:display_id=id, UGC:display_id=主播的bilin_id | Y            |              |
| avatar     | 主播头像                                          |              | Y            |
| start_at   | 开播时刻                                          |              | Y            |
| user_num   | 房间人数                                          |              | Y            |
| tag_url    | 房间标签（string array)                           |              | Y            |

## 同时搜索user和room

使用 `typ=-1` 搜索

## song

使用 `typ=3` 搜索，使用 topic `bilin_song_update` 或 `bilin_song_update_test` 更新

字段名（区分大小写）, 字段的取值，类型统一，都是 string

| 字段名           | 字段的意义                      | 是否加入搜索 | 是否外部存储 |
| ---------------- | ------------------------------- | ------------ | ------------ |
| id               | song id                         |              |              |
| name             | 歌曲名                          | Y            |              |
| artist           | 歌手名                          | Y            |              |
| duration         | 歌曲时长                        |              | Y            |
| upload_by        | 上传者                          |              | Y            |
| lyric            | 歌词地址                        |              | Y            |
| lyric_md5        | 歌词md5                         |              | Y            |
| lyric_len        | 歌词大小                        |              | Y            |
| audio            | 歌曲地址                        |              | Y            |
| audio_md5        | 歌曲md5                         |              | Y            |
| audio_len        | 歌曲大小                        |              | Y            |
| instrumental     | 伴奏曲地址                      |              | Y            |
| instrumental_md5 | 伴奏曲md5                       |              | Y            |
| instrumental_len | 伴奏曲大小                      |              | Y            |
| pkg              | 歌曲、伴奏曲、和歌词的zip包地址 |              | Y            |
| pkg_md5          | zip包md5                        |              | Y            |
| pkg_len          | zip包大小                       |              | Y            |

* zip包的格式，zip包解开后有三个文件：

  * `1.lrc`
  * `2.mp3`
  * `3.mp3`

  其中，1是歌词，2是歌曲，3是伴奏。
