# 介绍
## my-qqbot

进入[后go-cqhttp时代](https://github.com/Mrs4s/go-cqhttp/issues/2471)一段时间后的兴趣使然之作，较为简单，本来就是为平时的上网冲浪添加一个更方便的窗口，目前实现的功能有：
```
命令列表：
1.帮助
2.登录哔哩哔哩
3.ai对话
4.订阅直播间
5.订阅动态
6.重置对话
7.取消订阅动态
8.每日新闻
9.取消每日新闻
```
### 功能介绍
1. 帮助：“/help”或 “/帮助”输出以上命令列表，较为敷衍的实现，未排序
2. 登录哔哩哔哩：“/登录哔哩哔哩”扫码登录成功后自动获取cookie和refresh_token，并保存在配置文件中，后续订阅直播间或动态时需要登录
3. ai对话：私聊不加“/”前缀则直接开始ai对话，群聊需要at机器人开始对话
4. 订阅直播间：“/订阅直播间”+“房间号”订阅指定bilibili直播间开播提醒
5. 订阅动态：“/订阅动态”+“用户名”订阅指定bilibili用户动态，需要登录的账号关注该用户
6. 重置对话：“/重置对话”重置当前ai对话
7. 取消订阅动态：“/取消订阅动态”取消订阅指定bilibili用户动态
8. 每日新闻：“/每日新闻”订阅每日新闻，每天晚上8点（其实是UTC时间的中午12点）发送b站热搜
9. 取消每日新闻：取消订阅新闻

# 部署
## 手动部署
首先需要一个在线qq机器人来与QQ直接交互，请参考[LLOneBot](https://github.com/LLOneBot/LLOneBot)或[NapCatQQ](https://github.com/NapNeko/NapCatQQ)
使用正向WS(websocket)连接，开放端口号为`3001`，免得选择困难症
在程序所在目录下创建一个`config.yaml`,内容参考项目的`config.example.yaml`
```yaml
self: 123456789 # 机器人QQ号
nickname: [机器人昵称] # 机器人昵称，暂时没有实现相关功能
bilibili: # 这一字段使用“登录哔哩哔哩”自动填写
  cookie: "" # 哔哩哔哩cookie
  refresh_token: "" # 哔哩哔哩refresh_token
ws: # ws连接配置
  address: ws://127.0.0.1:3001/ # 机器人程序ws地址
  token: "" # ws连接可能所需的鉴权token

ai_chat: # 目前推荐使用deep-seek的api，也可以使用openai的api（可能需要代理）
  key: "" # api鉴权的api key
  base_url: "https://api.deepseek.com/v1"  # 实现了openai接口的ai对话api地址，
  model: "deepseek-chat" # 调用的模型名
```
## docker部署
参考项目的`docker-compose.yaml`文件，使用docker-compose一键部署，需要把`config.yaml`中ws的地址改为`ws://napcat:3001/`
也可参照其内容用docker部署


