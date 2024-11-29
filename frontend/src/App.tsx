import {useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {EventsOn} from "../wailsjs/runtime/runtime"

interface CardInfo {
    cardNumber: string
    name: string
    idNumber: string
    birthday: string
    sex: string
    cardDate: string
    isCardExist: boolean
}

function App() {
    const [cardInfo, setCardInfo] = useState<CardInfo | null>(null);

    useEffect(() => {
        EventsOn("card-status", (info: CardInfo) => {
            setCardInfo(info)
        })
    }, []);

    if(!cardInfo?.isCardExist){
        return (
            <>
                <img src={logo} alt="logo" id='logo' />
                <h1>未插入健保卡</h1>
            </>

        )
    }
    else return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id='result' className='result'>
                <h2>健保卡資料</h2>
                <p>卡號：{cardInfo.cardNumber}</p>
                <p>姓名：{cardInfo.name}</p>
                <p>身分證字號：{cardInfo.idNumber}</p>
                <p>出生日期：{cardInfo.birthday}</p>
                <p>性別：{cardInfo.sex}</p>
                <p>發卡日期：{cardInfo.cardDate}</p>
            </div>
        </div>
    )
}

export default App
