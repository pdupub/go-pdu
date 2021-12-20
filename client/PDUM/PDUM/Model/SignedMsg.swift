//
//  SignedMsg.swift
//  PDUM
//
//  Created by Liu Peng on 2021/8/31.
//

import Foundation

struct SignedMsg: Hashable, Codable, Identifiable {
    var content: String
    var refs: [String]
    var signature: String
    var id: String {
        signature
    }
    
    var quantum: Quantum? {
        let data = Data(base64Encoded: content)
        if let resObject = (try? JSONSerialization.jsonObject(with: data!)) as? [String: Any]{
            let t = resObject["t"] as! Int
            let v = resObject["v"] as! Int
            let d = resObject["d"] as! String
            return Quantum(t:t,v:v,d:d)
        }
        return nil
    }
    
    var qData: QData? {
        let data = Data(base64Encoded: quantum!.d)
        if let resObject = (try? JSONSerialization.jsonObject(with: data!)) as? [String: Any]{
            let text = resObject["text"] as! String
            let quote =  resObject["quote"] is NSNull ? nil : resObject["quote"] as? String
            let resources = resObject["resources"] as! [[String:Any]]
            var resArray = [QRes]()
            resources.forEach { item in
                let ft = item["format"] as! Int
                let d = item["data"] is NSNull ? nil : Data(base64Encoded:item["data"] as! String)
                let url = item["url"] as! String
                let cs = item["cs"] as! String
                let pir = QRes(format: ft, data: d, url: url, cs: cs)
                resArray.append(pir)
            }
            return QData(text: text, quote: quote, resources: resArray)
        }
        return nil
    }
    
    var pBorn: PBorn? = nil
    var qProfile: QProfile? = nil
    
    struct Quantum : Hashable, Codable {
        var t: Int
        var v: Int
        var d: String
    }
    
    struct QRes : Hashable, Codable {
        var format: Int
        var data: Data?
        var url: String
        var cs: String
    }
    
    struct QData : Hashable, Codable {
        var text: String
        var quote: String?
        var resources: [QRes]
    }
    
    struct PBorn : Hashable, Codable {
        var addr: String
        var sigs: [String]
    }
    
    struct QProfile : Hashable, Codable {
        var name: String
        var email: String
        var bio: String
        var url: String
        var location: String
        var avatar: QRes
        var extra: String
    }
}
