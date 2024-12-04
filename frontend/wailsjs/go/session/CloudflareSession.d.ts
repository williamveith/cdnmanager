// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {models} from '../models';
import {cloudflare} from '../models';

export function DeleteKeyValue(arg1:string):Promise<void>;

export function DeleteKeyValues(arg1:Array<string>):Promise<void>;

export function GetAllEntries():Promise<Array<models.Entry>>;

export function GetAllEntriesFromKeys(arg1:Array<cloudflare.StorageKey>):Promise<Array<models.Entry>>;

export function GetAllKeys():Promise<Array<cloudflare.StorageKey>>;

export function GetAllValues():Promise<Array<string>>;

export function GetValue(arg1:string):Promise<string>;

export function InsertKVEntry(arg1:string,arg2:string,arg3:string):Promise<cloudflare.Response>;

export function Size():Promise<number|Array<cloudflare.StorageKey>>;

export function WriteEntries(arg1:Array<models.Entry>):Promise<void>;

export function WriteEntry(arg1:models.Entry):Promise<cloudflare.Response>;
