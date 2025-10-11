# 🧪 快速测试指南 - v2.6.0-simplified

**目的**: 快速验证智能简化UI的效果

---

## 🚀 **立即测试（1分钟）**

```bash
cd /root/apple-music-downloader
./apple-music-downloader-v2.6.0-simplified <your_music_url>
```

---

## ✅ **验证要点**

### **1. 固定位置刷新**
- ✅ 7首歌的Track列表固定在原位
- ✅ 只有百分比和速度在变化
- ✅ **没有向下滚动**

### **2. 无行覆盖**
- ✅ 每一行都清晰完整
- ✅ Track编号不会被覆盖
- ✅ 状态信息不会被截断

### **3. 单行显示**
- ✅ 每个Track占用1行
- ✅ 没有因歌名过长而换行

---

## 🔍 **终端宽度测试**

### **测试1：宽屏（120+字符）**
```bash
# 拖拽终端窗口变宽
# 观察：显示完整模式
# [3/7] Spanish Key (feat. Wayne Shorter, ...) (16bit/44.1kHz) ↓54% 4.4MB/s
```

### **测试2：中等（80-119字符）**
```bash
# 调整终端到中等宽度
# 观察：自动切换紧凑模式
# [3/7] Spanish Key (feat. ...) 16/44 ↓54% 4.4MB/s
```

### **测试3：窄屏（<80字符）**
```bash
# 调整终端到较窄
# 观察：自动切换极简模式
# [3/7] Spanish Key... ↓54%
```

---

## 📊 **对比旧版本**

### **旧版本（v2.6.0-fixed2）**：
```
Track 3 of 7: Spanish Key (feat. Wayne Shorter, Bennie Maupin, John McLaughlin, Chick Corea, ... (16bit/44.1kHz) - 下载中 54% (4.4 MB/s)ck 7 of 7: Feio ... - 下载中 97%
              ↑ 过长导致换行 ↑                                                                                              ↑ 被覆盖 ↑
```

### **新版本（v2.6.0-simplified）**：
```
[3/7] Spanish Key (feat. Wayne Shorter, ...) (16bit/44.1kHz) ↓54% 4.4MB/s
                                                                           ↑ 单行显示，永不换行 ↑
```

---

## 🎯 **预期效果**

### **正常运行**：
```
[1/7] Pharaoh's Dance (feat. ...) 16/44 ↓52% 4.4MB/s
[2/7] Bitches Brew (feat. ...) 16/44 ↓37% 4.2MB/s
[3/7] Spanish Key (feat. ...) 16/44 ↓54% 4.4MB/s
[4/7] John McLaughlin (feat. ...) 16/44 ⏸
[5/7] Miles Runs the Voodoo Down... 16/44 ↓66% 4.9MB/s
[6/7] Sanctuary (feat. ...) 16/44 ⏸
[7/7] Feio (feat. ...) 16/44 ↓97% 1.2MB/s
```

- 状态实时更新
- 固定位置刷新
- 无滚动，无覆盖

---

## 🐛 **如遇问题**

### **问题1：仍然有滚动**
**原因**: 可能使用了旧版本  
**解决**: 
```bash
./apple-music-downloader-v2.6.0-simplified --version
# 确认版本号包含 "simplified"
```

### **问题2：歌名显示太短**
**原因**: 终端太窄（<80字符）  
**解决**: 拖拽终端窗口变宽，或接受极简模式

### **问题3：特殊字符显示异常**
**原因**: 终端不支持Unicode emoji  
**解决**: 这是预期的，不影响功能

---

## 📞 **报告结果**

测试后，请告知：
- ✅ 完全正常 - 太好了！
- ⚠️ 有小问题 - 描述具体现象
- ❌ 仍有滚动 - 提供终端宽度（`tput cols`）

---

**祝测试顺利！** 🚀

