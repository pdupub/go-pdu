//
//  ContentView.swift
//  PDUM
//
//  Created by Liu Peng on 2021/8/16.
//

import SwiftUI

struct ContentView: View {
    @EnvironmentObject var modelData: ModelData

    @State private var txt: String = ""
    @State private var exInfo: String = ""
    @State private var address: String = "0xAF040ed5498F9808550402ebB6C193E2a73b860a"
    @State private var getUrl: String = "http://127.0.0.1:1323"
    @State private var privKey: String = "689ac13dc3f424c8d5a6ef07a2e443311fc40ae4c370dac127bf5c1267e1ac98"
    @State private var message: String = "Hello World!!"
    @State private var reference: String = ""
    @State private var passwd: String = ""
    
    public enum HTTPMethod : String {
        case get   = "GET"
        case post  = "POST"
        case put   = "PUT"
        case patch = "PATCH"
        case delete = "DELETE"
    }
        
    func request(_ url: String,  _ httpMethod: HTTPMethod ,parameters: [String: String], completion: @escaping ([String: Any]?, Error?) -> Void) {
        var components = URLComponents(string: url)!
        components.queryItems = parameters.map { (key, value) in
            URLQueryItem(name: key, value: value)
        }
        components.percentEncodedQuery = components.percentEncodedQuery?.replacingOccurrences(of: "+", with: "%2B")
        var request = URLRequest(url: components.url!)
        request.httpMethod = httpMethod.rawValue
        if httpMethod == HTTPMethod.post {
            request.addValue("application/json", forHTTPHeaderField:"Content-Type")
            request.addValue("application/json", forHTTPHeaderField:"Accept")
            request.httpBody = Data(self.txt.utf8)
        }
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            guard
                let data = data,                              // is there data
                let response = response as? HTTPURLResponse,  // is there HTTP response
                200 ..< 300 ~= response.statusCode,           // is statusCode 2XX
                error == nil                                  // was there no error
            else {
                completion(nil, error)
                return
            }
            
            let responseObject = (try? JSONSerialization.jsonObject(with: data)) as? [String: Any]
            completion(responseObject, nil)
        }
        task.resume()
    }

    var body: some View {
        VStack{
            Text(txt)
                .lineLimit(nil)
                .multilineTextAlignment(.leading)
            Spacer()
            HStack{
                Text("Server URL")
                TextField("", text: $getUrl)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            }
            HStack{
                Text("Private Key")
                TextField("", text: $privKey)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            }
            HStack{
                Text("ExtraInfo")
                TextField("", text: $exInfo)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            }
            HStack{
                Text("Message")
                TextField("", text: $message)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            }
            HStack{
                Text("Reference")
                TextField("", text: $reference)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            }
            HStack{
                Text("Password")
                TextField("", text: $passwd)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            }

            HStack{
                Spacer()
                Button("Reference"){
                    request(getUrl+"/info/latest/0xAF040ed5498F9808550402ebB6C193E2a73b860a",HTTPMethod.get, parameters: ["foo": "bar"]) { responseObject, error in
                        guard let responseObject = responseObject, error == nil else {
                            print(error ?? "Unknown error")
                            return
                        }
                        self.reference = responseObject["signature"] as! String
                    }
                }
                Button("SignMsg"){
                    // Reverse text here
                    let str = signMsg(UnsafeMutablePointer<Int8>(mutating: (self.privKey as NSString).utf8String),UnsafeMutablePointer<Int8>(mutating: (self.message as NSString).utf8String),UnsafeMutablePointer<Int8>(mutating: (self.reference as NSString).utf8String))
                    self.txt = String.init(cString: str!, encoding: .utf8)!
                    // don't forget to release the memory to the C String
                    str?.deallocate()
                    
                    request(getUrl, HTTPMethod.post, parameters:[:]){ responseObject, error in
                        guard let responseObject = responseObject, error == nil else {
                            print(error ?? "Unknown error")
                            return
                        }
                        self.reference = responseObject["signature"] as! String
                    }
                }
                Button("ExtraInfo") {
//                    self.exInfo = modelData.profile.username
//                    self.exInfo = modelData.smsgs[0].signature
//                    self.exInfo = modelData.smsgs[0].quantum!.d
                    if (modelData.smsgs[0].qdata!.resources[0] == 1) {
                        self.exInfo = modelData.smsgs[0].qdata!.resources[0].url
                    }
                }
                Button("GetAddress") {
                    // Reverse text here
                    let str = getAddress(UnsafeMutablePointer<Int8>(mutating: (self.privKey as NSString).utf8String))
                    self.txt = String.init(cString: str!, encoding: .utf8)!
                    // don't forget to release the memory to the C String
                    str?.deallocate()
                }
                Button("GenerateKey") {
                    // Reverse text here
                    let str = generateKey()
                    self.txt = String.init(cString: str!, encoding: .utf8)!
                    self.privKey = self.txt
                    // don't forget to release the memory to the C String
                    str?.deallocate()
                }
                Button("GenerateKeystore") {
                    // Reverse text here
                    let str = generateKeystore(UnsafeMutablePointer<Int8>(mutating: (self.privKey as NSString).utf8String), UnsafeMutablePointer<Int8>(mutating: (self.passwd as NSString).utf8String))
                    self.txt = String.init(cString: str!, encoding: .utf8)!
                    // don't forget to release the memory to the C String
                    str?.deallocate()
                }
                Button("UnlockKeystore") {
                    // Reverse text here
                    let str = unlockKeystore(UnsafeMutablePointer<Int8>(mutating: (self.txt as NSString).utf8String), UnsafeMutablePointer<Int8>(mutating: (self.passwd as NSString).utf8String))
                    self.txt = String.init(cString: str!, encoding: .utf8)!
                    // don't forget to release the memory to the C String
                    str?.deallocate()
                }
                Button("Ecrecover") {
                    // Reverse text here
                    let str = ecrecover(UnsafeMutablePointer<Int8>(mutating: (self.txt as NSString).utf8String))
                    self.txt = String.init(cString: str!, encoding: .utf8)!
                    // don't forget to release the memory to the C String
                    str?.deallocate()
                }
                

                Spacer()
            }
            
            Spacer()
        }
        .padding(.all, 15)
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView().environmentObject(ModelData())
    }
}
