'use strict';

const Boom = require('@hapi/boom');
const GetConversion = require('../../application/use_cases/GetConversion');
const ListPrices = require('../../application/use_cases/ListPrices');

module.exports = {

  async getListDataPrice(request) {

    // Context
    const serviceLocator = request.server.app.serviceLocator;

    // Treatment
    try {
      const prices = await ListPrices(serviceLocator);
      // Output
      return {"prices" : prices, "total" : prices.length};
    } catch (err) {
      return Boom.unauthorized('Bad credentials');
    }
  },

  async getAggregatePrice(request) {
    // Context
    const serviceLocator = request.server.app.serviceLocator;
    // Input
    const dateFrom = request.query.dateFrom
    const dateTo = request.query.dateTo
    const areaProvinsi = request.query.areaProvinsi
    // Treatment
    try {
      const aggregateResult = await GetConversion(dateFrom, dateTo, areaProvinsi, serviceLocator);
      // Output
      return aggregateResult;
    } catch (err) {
      console.log(err)
      return Boom.unauthorized('Bad credentials');
    }
  },

};