'use strict';

const sequelize = require('../orm/sequelize/sequelize');
const User = require('../../domain/User');
const UserRepository = require('../../domain/UserRepository');

module.exports = class extends UserRepository {

  constructor() {
    super();
    this.db = sequelize;
    this.model = this.db.model('user');
  }

  async persist(userEntity) {
    const {  roleId, name, phone, password } = userEntity;
    const seqUser = await this.model.create({  roleId, name, phone, password});
    await seqUser.save();

    return new User(seqUser.id, seqUser.roleId, seqUser.name, seqUser.phone, seqUser.password);
  }

  async merge(userEntity) {
    const seqUser = await this.model.findByPk(userEntity.id);

    if (!seqUser) return false;

    const { roleId, name, phone, password, token, tokenExpiredAt, updatedAt, updatedBy, deletedAt, deletedBy } = userEntity;
    await seqUser.update({ roleId, name, phone, password, token, tokenExpiredAt, updatedAt, updatedBy, deletedAt, deletedBy});

    return new User(seqUser.id, seqUser.roleId, seqUser.name, seqUser.phone, seqUser.password);
  }

  async remove(userId) {
    const seqUser = await this.model.findByPk(userId);
    if (seqUser) {
      return seqUser.destroy();
    }
    return false;
  }

  async get(userId) {
    const seqUser = await this.model.findByPk(userId);
    return new User(seqUser.id, seqUser.roleId, seqUser.name, seqUser.phone, seqUser.password);
  }

  async getByPhone(userPhone) {
    const seqUser = await this.model.findOne({ where: { phone: userPhone } });
    return new User(seqUser.id, seqUser.roleId, seqUser.name, seqUser.phone, seqUser.password);
  }

  async find() {
    const seqUsers = await this.model.findAll();
    return seqUsers.map((seqUser) => {
      return new User(seqUser.id, seqUser.roleId, seqUser.name, seqUser.phone, seqUser.password);
    });
  }

};
