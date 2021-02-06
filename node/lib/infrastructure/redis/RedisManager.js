'use strict';

const RedisManager = require('../../application/redis/RedisManager');

var redis = require('redis');
const { promisify } = require('util');

const port = "6379"
const host = "127.0.0.1"

var client = redis.createClient(port, host); //creates a new client

client.on('connect', function(){
  console.log('Redis Connected !')
})

var date = new Date()


module.exports = class extends RedisManager {
  
  setData(key, value) {
    client.set(key, value)
    client.expire(key, date.getHours() * 12)
  }

  async getData(key) {
    const getAsync = promisify(client.get).bind(client)
    const value = await getAsync(key)
    return value
  }

  removeData(key) {
    client.del(key)
  }

};