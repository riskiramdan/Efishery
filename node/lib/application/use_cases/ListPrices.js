'use strict';

const fetch = require('node-fetch');

module.exports = async ( { redisManager }) => {  
  const responsePrices = await fetch("https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list", {method:'Get'})
  const bodyPrices = await responsePrices.json()

  var USDPrice = await redisManager.getData("USD_IDR")
  if (USDPrice == null) {
    const responseUSD = await fetch("https://free.currconv.com/api/v7/convert?q=USD_IDR&compact=ultra&apiKey=bbc99fbf718c1a23be93", {method:'Get'})
    const bodyUSD = await responseUSD.json()
    USDPrice = bodyUSD.USD_IDR
    redisManager.setData("USD_IDR", bodyUSD.USD_IDR)
  }
  
  var arr = []
  for(var i = 0;i<bodyPrices.length;i++){
    if (bodyPrices[i].uuid != null) {
      var a = {
        "uuid" : bodyPrices[i].uuid,
        "komoditas" : bodyPrices[i].komoditas,
        "area_provinsi" : bodyPrices[i].area_provinsi,
        "area_kota" : bodyPrices[i].area_kota,
        "size" :bodyPrices[i].size,
        "price" : parseFloat(bodyPrices[i].price),
        "priceUSD" : parseFloat(bodyPrices[i].price) / parseFloat(USDPrice),
        "tgl_parsed" :bodyPrices[i].tgl_parsed,
        "timestamp":bodyPrices[i].timestamp
      }
      arr.push(a)
    }
  }
  
  return arr
};

