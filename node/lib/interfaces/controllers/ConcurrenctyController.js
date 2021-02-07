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

  async getConversionPrice(request) {
    // Context
    const serviceLocator = request.server.app.serviceLocator;
    // Input
    const dateFrom = request.query.dateFrom
    const dateTo = request.query.dateTo
    const areaProvinsi = request.query.areaProvinsi
    // Treatment
    try {
      const conversion = await GetConversion(dateFrom, dateTo, areaProvinsi, serviceLocator);
      // Output
      return conversion;
    } catch (err) {
      console.log(err)
      return Boom.unauthorized('Bad credentials');
    }
  },

};