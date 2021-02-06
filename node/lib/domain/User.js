'use strict';

module.exports = class {

  constructor(id = null, roleId, name, phone, password, token, tokenExpiredAt, createdAt, createdBy, updatedAt, updatedBy, deletedAt, deletedBy) {
    this.id = id;
    this.roleId = roleId;
    this.name = name,
    this.phone = phone,
    this.password = password,
    this.token = token,
    this.tokenExpiredAt = tokenExpiredAt,
    this.createdAt = createdAt,
    this.createdBy = createdBy,
    this.updatedAt = updatedAt,
    this.updatedBy = updatedBy,
    this.deletedAt = deletedAt,
    this.deletedBy = deletedBy
  }

};