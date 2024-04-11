import {api_url} from "./index.js"
import { get } from "svelte/store"

export class Article {
  /**
   * @param {string} id 
   * @param {string} link 
   * @param {string} title 
   * @param {Date} pubdate 
   * @param {string} pubid 
   * @param {string} summary 
   */
  constructor(id, link, title, pubdate, pubid, summary) {
    this.id = id
    this.link = link
    this.title = title
    this.pubdate = pubdate
    this.pubid = pubid
    this.summary = summary
  }
}
/**
 * @param {number} limit
 * @param {number} page 
 * @param {string} pub
 * @param {boolean} sum
 * @param {boolean} audio
 * @returns {Promise<Article[]>}
 */
export async function articleFromApi(limit, page, pub, sum, audio) {
  /**
   * @type {Array<any>}
   */
  let jsResponse = await fetch(get(api_url)+
    `/articles?limit=${limit}&start=${page*limit}&PublisherID=${pub}&summary=${sum}&audio=${audio}`)
    .then(res => res.json())
    .catch(err => {console.log(err);return []})
  /**
   * @type {Article[]}
   */
  let articles = []
  jsResponse.map(v => {
    articles.push(new Article(
      v["ArticleID"], v["Link"], v["Title"], v["PubDate"], v["PublisherID"], v["Summary"]
    ))
  })
  if (articles.length === 0 && page > 0) {
    location.reload()
    return []
  }
  return articles
}