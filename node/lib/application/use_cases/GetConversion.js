'use strict';

const ListPrices = require('./ListPrices');
const jsonQ = require('js-jsonq');

module.exports = async (dateFrom, dateTo, areaProvinsi, { redisManager }) => {
  var listPrices = await redisManager.getData("LISTPRICES")
  if (listPrices == null) {
    listPrices = await ListPrices( {redisManager} )
  }else{
    listPrices = JSON.parse(listPrices)
  }

  var query = new jsonQ(listPrices)
  var from, to
  if (dateFrom != undefined) {
    from = Date.parse(dateFrom)
    query = query.where("timestamp", ">", from)   
  }
  if (dateTo != undefined) {
    to = Date.parse(dateTo)
    query = query.where("timestamp", "<", to)
  }
  if (areaProvinsi !== undefined) {
    query = query.where("area_provinsi", "contains", areaProvinsi)
  }
  return {
    "price" : {
      "min" : query.min("price"),
      "max" : query.max("price"),
      "median" : query.sortBy("price", "asc").sum("price") / 2,
      "average" : query.avg("price"),
    },"total" : query.count()
  }
};

