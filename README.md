# Go-NHI-Reader-Wails

## About

類似於[Benknightdark/ElectronVueCardReader](https://github.com/Benknightdark/ElectronVueCardReader)，但是本專案使用Wails/React構建。

這個專案可供讀取健保IC卡內資料，並列出以下資訊：

- 健保卡卡號
- 姓名
- 身分證字號
- 出生年月日
- 性別
- 發卡日期


## Usage

開始執行或測試之前，您需要：
- Go 1.23.2(或以上)
- [Wails](https://wails.io/)
- node.js 20.x LTS
- pnpm@latest
- 一個功能正常的讀卡機(且已有完備的驅動程式)

執行步驟如下：

1. 進入./frontend，執行`pnpm i`。

2. 回到專案根目錄，執行`wails dev`。

- 若您想編譯為獨立程式，請執行`wails build`。